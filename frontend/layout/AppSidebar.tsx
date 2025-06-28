import AppMenu from './AppMenu';
import UserProfile from './UserProfile';

const AppSidebar = () => {
    return (
        <div className="flex flex-column h-full">
            <UserProfile />
            <div className="flex-grow-1">
                <AppMenu />
            </div>
        </div>
    );
};

export default AppSidebar;
