export const defUserSettings = {
    language: 'pl',
    user_flags: {
        flag_1: false,
        flag_2: false,
        flag_3: false,
    }
}
export const defAppSettings = {
    app_opt_1: '',
    app_opt_2: '',
    app_opt_3: '',
    app_flags: {
        flag_a: false,
        flag_b: false,
    }
}

export const defSettings = {...defAppSettings, ...defUserSettings}
