---
trigger: always_on
---

# **Tools** - When to use which tools

## Browser access
If you need to test the results of your work, open your browser and test the results of your work.
Check the local environment url by running this command line tool in terminal:
`./.agent/scripts/get_domain.sh`
If you need to call api directly, add `/api/{handling_endpoint}` to the application url.
The email and password for this environment are: testagent@test.pl/t3$T4gent127

## Mysql query test
Use the provided helper script to query the database in your local environment. You will be automatically logged in when using this tool. You can check the database structure or entries when testing local changes.

To execute an SQL query, use the appropriate command line tool:

```bash
./.agent/scripts/run_db_query.sh <sql_query>
```

Parameters:
- **sql_query**: the SQL query to execute

Examples:

```bash
./.agent/scripts/run_db_query.sh "SELECT COUNT(*) FROM users"
./.agent/scripts/run_db_query.sh "SELECT * FROM users LIMIT 5"
```

### Important Notes
- **NEVER** try to connect to the database directly using `mysql` or `mariadb` from the host terminal.
- **NEVER** ask the user for database credentials. The script handles them.
- If you need to verify the result of a previous operation (e.g. "did the user get created?"), use this script to check.
