export type UserSettingsType = {
    language: string
    user_flags:Record<string, boolean>
}
export type AppSettingsType = {
    app_flags:Record<string, boolean>
    app_opt_1: string
    app_opt_2: string
    app_opt_3: string
}
export type SettingsType = UserSettingsType & AppSettingsType