/* eslint-disable @next/next/no-img-element */
'use client';
import { useRouter } from 'next/navigation';
import React, { useContext, useEffect, useRef, useState } from 'react';
import { Checkbox } from 'primereact/checkbox';
import { Button } from 'primereact/button';
import { Password } from 'primereact/password';
import { LayoutContext } from '../../../../layout/context/layoutcontext';
import { InputText } from 'primereact/inputtext';
import { classNames } from 'primereact/utils';
import { Toast } from 'primereact/toast';
import { authApi } from '@/app/services/auth';
import { useAuth } from '@/hooks/useAuth';

const LoginPage = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [checked, setChecked] = useState(false);
    const { layoutConfig } = useContext(LayoutContext);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const toast = useRef<Toast>(null);
    const router = useRouter();
    const isAuthenticated = useAuth();

    useEffect(() => {
        if (isAuthenticated) {
            router.push('/');
        }
    }, [isAuthenticated, router]);

    const containerClassName = classNames('surface-ground flex align-items-center justify-content-center min-h-screen min-w-screen overflow-hidden', { 'p-input-filled': layoutConfig.inputStyle === 'filled' });

    const showToast = (severity: 'success' | 'error', summary: string, detail: string) => {
        toast.current?.show({ severity, summary, detail, life: 3000 });
    };

    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError('');

        try {
            await authApi.login(email, password);
            router.push('/');
        } catch (error: any) {
            if (toast.current) {
                error.messages?.forEach((msg: string) => {
                    showToast('error', `Erro ${error.code}`, msg);
                });
            }
        } finally {
            setLoading(false);
        }
    };

    return (
        <>
            <Toast ref={toast} />
            <div className={containerClassName}>
                <div className="flex flex-column align-items-center justify-content-center">
                    <img src={`/layout/images/logo_eskcti.png`} alt="Eskokado Consultoria de TI logo" className="mb-5 w-6rem flex-shrink-0" />
                    <div
                        style={{
                            borderRadius: '56px',
                            padding: '0.3rem',
                            background: 'rgba(38, 150, 243, 0)'
                        }}
                    >
                        <form onSubmit={handleLogin}>
                            <div className="w-full surface-card py-8 px-5 sm:px-8" style={{ borderRadius: '53px' }}>
                                <div>
                                    <label htmlFor="email1" className="block text-900 text-xl font-medium mb-2">
                                        Email
                                    </label>
                                    <InputText
                                        id="email1" type="email" placeholder="Email address"
                                        className="w-full md:w-30rem mb-5"
                                        style={{ padding: '1rem' }}
                                        value={email}
                                        onChange={(e) => setEmail(e.target.value)}
                                    />

                                    <label htmlFor="password1" className="block text-900 font-medium text-xl mb-2">
                                        Password
                                    </label>
                                    <Password
                                        inputId="password1"
                                        value={password}
                                        onChange={(e) => setPassword(e.target.value)}
                                        placeholder="Password"
                                        toggleMask className="w-full mb-5"
                                        inputClassName="w-full p-3 md:w-30rem">
                                    </Password>

                                    <div className="flex align-items-center justify-content-between mb-5 gap-5">
                                        <a
                                            className="font-medium no-underline ml-2 text-right cursor-pointer"
                                            style={{ color: 'var(--primary-color)' }}
                                            onClick={() => router.push('/auth/register')}
                                        >
                                            Registrar
                                        </a>
                                        <a
                                            className="font-medium no-underline ml-2 text-right cursor-pointer"
                                            style={{ color: 'var(--primary-color)' }}
                                            onClick={() => router.push('/auth/forgot_password')}
                                        >
                                            Esqueci a senha
                                        </a>
                                    </div>

                                    {error && <small className="p-error">{error}</small>}

                                    <Button
                                        label="Entrar"
                                        className="w-full p-3 text-xl"
                                        type="submit"
                                        loading={loading}
                                    />
                                </div>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </>
    );
};

export default LoginPage;
