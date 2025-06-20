'use client';
import { LayoutProvider } from '../layout/context/layoutcontext';
import { PrimeReactProvider } from 'primereact/api';
import 'primereact/resources/primereact.css';
import 'primeflex/primeflex.css';
import 'primeicons/primeicons.css';
import '../styles/layout/layout.scss';
import '../styles/demo/Demos.scss';
import { CurrentUserProvider, useCurrentUser, UserDataType } from '@/layout/context/CurrentUserContext';
import { useEffect } from 'react';

interface RootLayoutProps {
    children: React.ReactNode;
}

export default function RootLayout({ children }: RootLayoutProps) {
    return (
        <html lang="en" suppressHydrationWarning>
            <head>
                <link id="theme-css" href={`/themes/lara-light-indigo/theme.css`} rel="stylesheet"></link>
            </head>
            <body>
                <PrimeReactProvider>
                    <CurrentUserProvider>
                        <LayoutProvider>
                            <InitCurrentUserContext />
                            {children}
                        </LayoutProvider>
                    </CurrentUserProvider>
                </PrimeReactProvider>
            </body>
        </html>
    );
}


const InitCurrentUserContext = () => {
    const { setCurrentUser } = useCurrentUser();

    useEffect(() => {
        let storeCurrentUser: UserDataType = {
            id: '',
            name: '',
            email: ''
        };
        storeCurrentUser.id = localStorage.getItem('user-id') || '';
        storeCurrentUser.name = localStorage.getItem('user-name') || '';
        storeCurrentUser.email = localStorage.getItem('user-email') || '';
        setCurrentUser(storeCurrentUser);
    }, [setCurrentUser]);

    return null;
};
