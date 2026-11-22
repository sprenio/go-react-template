export type LangType = {
  id: number;
  code: string;
  name: string;
};
export type ConfigType = {
  Features: Record<string, boolean>;
  Languages: LangType[];
  AppName: string;
  DefaultLanguage: string;
};
