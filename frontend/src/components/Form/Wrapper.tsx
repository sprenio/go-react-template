import React from 'react';

type Props = {
    children: React.ReactNode;
    onSubmit: (e: React.FormEvent<HTMLFormElement>) => void;
};

export const Wrapper = ({children, onSubmit}:Props) => {
    return (
        <form className={"form-wrapper"} onSubmit={onSubmit}>
            {children}
        </form>
    )
}