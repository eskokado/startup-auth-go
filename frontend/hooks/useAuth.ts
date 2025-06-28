import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';

export const useAuth = (redirectPath = '/auth/login') => {
    const router = useRouter();
    const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null);

    useEffect(() => {
        const checkAuth = () => {
            const token = localStorage.getItem('access-token');

            const authStatus = !!token;
            setIsAuthenticated(authStatus);

            if (!authStatus && redirectPath) {
                router.push(redirectPath);
            }
        };

        // Verificação no client-side
        if (typeof window !== 'undefined') {
            checkAuth();
        }
    }, [router, redirectPath]);

    return isAuthenticated;
};
