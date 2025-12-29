export type LangType = {
    id: number,
    code: string,
    i18n_code: string,
    name: string
}
export type ConfigType = {
    Features: Record<string, boolean>,
    Languages: LangType[],
    AppName: string,
    DefaultLanguage: string
}