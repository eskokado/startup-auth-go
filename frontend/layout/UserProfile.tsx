// components/layout/UserProfile.tsx
import React from 'react';
import { useCurrentUser } from './context/CurrentUserContext';

const UserProfile = () => {
    const { currentUser } = useCurrentUser();

    if (!currentUser || !currentUser.name) {
        return null;
    }

    return (
        <div className="user-profile p-3 border-bottom-1 surface-border">
            <div className="flex align-items-center">
                <img
                    src="/demo/images/avatar/amyelsner.png"
                    alt="User"
                    className="mr-3 w-3rem h-3rem border-circle"
                />
                <div className="flex flex-column">
                    <span className="font-bold">{currentUser.name}</span>
                    <span className="text-sm text-color-secondary">{currentUser.email}</span>
                </div>
            </div>
        </div>
    );
};

export default UserProfile;
