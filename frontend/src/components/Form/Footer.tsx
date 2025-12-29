import React from "react";

export const Footer = ({children}:{children: React.ReactNode}) => {
    return (
        <div className={"form-footer"}>
            {children}
        </div>
    )
}