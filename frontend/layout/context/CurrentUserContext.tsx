import React, { createContext, useContext, useState, ReactNode } from 'react';

interface CurrentUserContextType {
    currentUser: UserDataType | null;
    setCurrentUser: (user: UserDataType | null) => void;
}

export interface UserDataType {
    id: string | null;
    name: string | null;
    email: string | null;
}

const CurrentUserContext = createContext<CurrentUserContextType | undefined>(undefined);

export const CurrentUserProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
    const [currentUser, setCurrentUser] = useState<UserDataType | null>(null);

    return (
        <CurrentUserContext.Provider value={{
            currentUser,
            setCurrentUser,
        }}>
            {children}
        </CurrentUserContext.Provider>
    );
};

export const useCurrentUser = () => {
    const context = useContext(CurrentUserContext);
    if (context === undefined) {
        throw new Error('useCurrentUser must be used within an CurrentUserProvider');
    }
    return context;
};
