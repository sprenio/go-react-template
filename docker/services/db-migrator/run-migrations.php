<?php
class DbMigrator {
    private PDO $db;
    private bool $connected = false;

    public function __construct()
    {
        $this->setDbConnection();
    }

    public function doMigration(): bool
    {
        if(!$this->connected){
           $this->printMessage('Database connection failed');
           return false;
        }
        $sqlDir = dirname(__FILE__) . '/migrations';

        if (!file_exists($sqlDir)) {
            $this->printMessage('Migrations directory does not exist');
            return true; // this is no blocking error
        }
        $minDateDir = $this->getLastDateDir();
        $dateDirs = glob($sqlDir . DIRECTORY_SEPARATOR . '20????');
        sort($dateDirs);
        $currDecade = intdiv(date('y'), 10);
        $currDateDir = date('Ym');
        $dirPattern = '/^20[2-' . $currDecade . ']\\d(0\\d|1[0-2])$/';
        foreach ($dateDirs as $dateDir){
            if(!is_dir($dateDir)){
                continue;
            }
            $dirName = basename($dateDir);
            if(preg_match($dirPattern, $dirName)){
                if($dirName > $currDateDir){
                    continue;
                } elseif ($minDateDir > $dirName){
                    continue;
                }
            } else {
                continue;
            }
            $this->printMessage('Processing directory ' . $dateDir);

            $files = glob($dateDir . DIRECTORY_SEPARATOR . '*.sql');
            sort($files);
            foreach($files as $file){
                if(!is_file($file)){
                    continue;
                }
                $this->printMessage('Processing file ' . $file);
                if($this->fileHasBeenProcessed($file)){
                    $this->printMessage('file ' . $file . ' has been already processed');
                    continue;
                }
                $this->db->beginTransaction();
                $result = $this->processSqlFile($file);
                if(!$result){
                    if($this->db->inTransaction()){
                        $this->db->rollBack();
                    }
                    $this->printMessage('processing of files stopped due to errors');
                    return false;
                }
                if($this->db->inTransaction()){
                    $this->db->commit();
                }
            }
        }
        return true;
    }

    private function getLastDateDir(): int
    {
        $sql = 'SELECT MAX(`date_dir`) AS lastDateDir FROM `db_changes`';
        $stmt = $this->db->query($sql);
        $row = $stmt->fetch(PDO::FETCH_ASSOC);
        return $row['lastDateDir'] ?? 0;
    }

    /**
     * @param string $dsn
     * @return PDO
     * @throws PDOException
     */
    private function connectWithRetry(string $dsn, string $rootPass): PDO
    {
        $maxAttempts = 5;
        $attempt = 0;
        do {
            if($attempt > 0){
                $backoff = min(pow(2, $attempt), 30);
                $this->printMessage('Attempt ' . $attempt . ' failed. Retrying in ' . $backoff . ' seconds...');
                sleep($backoff);
            }
            $attempt++;
            try {
                $db = new PDO($dsn, 'root', $rootPass, [
                  PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION
                ]);
                return $db;
            } catch (PDOException $e) {
                echo "Attempt " . $attempt . " Error: " . $e->getMessage() . PHP_EOL;
            }
        } while ($attempt < $maxAttempts);
        throw new PDOException("Failed to connect after " . $maxAttempts . " attempts");
    }
    private function setDbConnection(): void {
        $dbHost = getenv('MYSQL_HOST');
        $dbPass = getenv('MYSQL_ROOT_PASSWORD');
        $dbName = getenv('MYSQL_DATABASE');
        if(empty($dbHost) || empty($dbPass) || empty($dbName)){
            $this->printMessage('Missing environment variables');
            return;
        }

        $dsn = 'mysql:host=' . $dbHost . ';dbname=' . $dbName . ';charset=utf8mb4';
        try {
            $this->db = $this->connectWithRetry($dsn, $dbPass);
        } catch (PDOException $e) {
            echo "Error: " . $e->getMessage() . PHP_EOL;
            return;
        }
        $this->connected = true;
    }
    private function printMessage(string $message): void
    {
        echo date('Y-m-d H:i:s - ') . $message . PHP_EOL;
    }
    private function getSqlFileInfo(string $file): array
    {
        $info = [
          'fileName' => basename($file),
          'dateInt' => 0,
        ];
        $dirName = basename(dirname($file));
        if(is_numeric($dirName)){
            $dirInt = intval($dirName);
            if($dirInt > 202508 && $dirInt <= intval(date('Ym'))){
                $info['dateInt'] = $dirInt;
            }
        }
        return $info;
    }

    private function fileHasBeenProcessed(string $file): bool
    {
        $info = $this->getSqlFileInfo($file);
        if(empty($info['dateInt'])){
            return false;
        }
        return $this->getDbChangeId($info['fileName'], $info['dateInt'], true) > 0;
    }
    private function getDbChangeId(string $fileName, int $dirInt, bool $onlyCompleted = false): int
    {
        $sql = 'SELECT id FROM `db_changes` WHERE `file_name` = ? AND `date_dir` = ?';
        if($onlyCompleted){
            $sql .=' AND `complete_date` IS NOT NULL';
        }
        $stmt = $this->db->prepare($sql);
        $stmt->execute([$fileName, $dirInt]);
        $row = $stmt->fetch(PDO::FETCH_ASSOC);
        return $row['id'] ?? 0;
    }
    private function processSqlFile(string $file): bool
    {
        if(!is_readable($file)){
            $this->printMessage('file ' . $file . ' is not readable');
            return false;
        }
        $info = $this->getSqlFileInfo($file);
        if(empty($info['dateInt'])){
            $this->printMessage('file ' . $file . ' has invalid date');
            return false;
        }
        $id = $this->getDbChangeId($info['fileName'], $info['dateInt']);
        if(empty($id)){
            $sql = 'INSERT INTO `db_changes` (`file_name`, `date_dir`, `start_date`) VALUES (?, ?, NOW())';
            $stmt = $this->db->prepare($sql);
            $stmt->execute([$info['fileName'], $info['dateInt']]);
            $id = $this->db->lastInsertId();
            if(empty($id)){
                $this->printMessage('invalid last insert id for file ' . $file);
                return false;
            }
        }
        $sql = file_get_contents($file);
        try {
            if(!empty($sql)) {
                $this->db->exec($sql);
            }
            $sql = 'UPDATE `db_changes` SET `complete_date` = NOW() WHERE `id` = ?';
            $stmt = $this->db->prepare($sql);
            $stmt->execute([$id]);
        } catch (PDOException $e) {
            $this->printMessage('file ' . $file . ' processing failed: ' . $e->getMessage());
            return false;
        }
        $this->printMessage('file ' . $file . ' processed successfully');
        return true;
    }
}

echo date('Y-m-d H:i:s') . " - Starting database migrations...\n";
$migrator = new DbMigrator();
if($migrator->doMigration()){
    echo date('Y-m-d H:i:s') . " - Database migrations completed successfully.\n";
} else {
    echo date('Y-m-d H:i:s') . " - Database migrations failed.\n";
    exit(1);
}
