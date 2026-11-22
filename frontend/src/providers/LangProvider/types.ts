export type LangContextType = {
  lang: string;
  locale: string;
  changeLang: (newLang: string) => void;
};
