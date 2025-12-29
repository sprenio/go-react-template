type LoaderState = {
    visible: boolean
    message: string
}

export type LoaderContextType = {
    loader: LoaderState
    showLoader: (message?: string) => void
    hideLoader: () => void
}
