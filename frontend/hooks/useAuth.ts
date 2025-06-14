// src/hooks/useAuth.ts
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';

export const useAuth = (redirectPath = '/auth/login') => {
    const router = useRouter();
    const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null);

    useEffect(() => {
        const checkAuth = () => {
            const token = localStorage.getItem('access-token');
            const client = localStorage.getItem('client');
            const uid = localStorage.getItem('uid');

            const authStatus = !!token && !!client && !!uid;
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
