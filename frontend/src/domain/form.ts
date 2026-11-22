export const FormModeEdit = 'edit' as const;
export const FormModeCreate = 'create' as const;
export const FormModeDelete = 'delete' as const;

export type FormMode = typeof FormModeCreate | typeof FormModeEdit | typeof FormModeDelete;
