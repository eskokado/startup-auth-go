/* eslint-disable @next/next/no-img-element */
'use client';
import { useRouter, useSearchParams } from 'next/navigation';
import React, { useContext, useEffect, useState } from 'react';
import { Button } from 'primereact/button';
import { Password } from 'primereact/password';
import { LayoutContext } from '../../../../layout/context/layoutcontext';
import { InputText } from 'primereact/inputtext';
import { classNames } from 'primereact/utils';
import { authApi } from '@/app/services/auth';

const UpdatePasswordPage = () => {
    const searchParams = useSearchParams();
    const [token, setToken] = useState<string>('');
    const [password, setPassword] = useState('');
    const [confirmationPassword, setConfirmationPassword] = useState('');
    const { layoutConfig } = useContext(LayoutContext);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    useEffect(() => {
        const urlToken = searchParams.get('reset_password_token');
        console.log("token", urlToken)
        setToken(urlToken ?? '');
    }, []);

    const router = useRouter();

    const containerClassName = classNames('surface-ground flex align-items-center justify-content-center min-h-screen min-w-screen overflow-hidden', { 'p-input-filled': layoutConfig.inputStyle === 'filled' });

    const handleUpdatePassword = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError('');

        try {
            if (password !== confirmationPassword) {
                throw new Error('Passwords do not match');
            }

            await authApi.updatePassword({
                reset_password_token: token,
                password
            });

            router.push('/auth/login');
        } catch (err: any) {
            setError(err.message || 'Ocorreu falha ao alterar senha');
        } finally {
            setLoading(false);
        }
    };

    return (
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
                    <form onSubmit={handleUpdatePassword}>
                        <div className="w-full surface-card py-8 px-5 sm:px-8" style={{ borderRadius: '53px' }}>
                            <div>
                                <label htmlFor="password" className="block text-900 font-medium text-xl mb-2">
                                    Password
                                </label>
                                <Password
                                    inputId="password"
                                    value={password}
                                    onChange={(e) => setPassword(e.target.value)}
                                    placeholder="Password"
                                    toggleMask
                                    className="w-full mb-5"
                                    inputClassName="w-full p-3 md:w-30rem"
                                />

                                <label htmlFor="confirmPassword" className="block text-900 font-medium text-xl mb-2">
                                    Confirm Password
                                </label>
                                <Password
                                    inputId="confirmPassword"
                                    value={confirmationPassword}
                                    onChange={(e) => setConfirmationPassword(e.target.value)}
                                    placeholder="Confirm Password"
                                    toggleMask
                                    className="w-full mb-5"
                                    inputClassName="w-full p-3 md:w-30rem"
                                    feedback={false}
                                />
                                <div className="flex align-items-center justify-content-between mb-5 gap-5">
                                    <a
                                        className="font-medium no-underline ml-2 text-right cursor-pointer"
                                        style={{ color: 'var(--primary-color)' }}
                                        onClick={() => router.push('/auth/login')}
                                    >
                                        Logar?
                                    </a>
                                </div>

                                {error && <small className="p-error block mb-5">{error}</small>}

                                <Button
                                    label="Alterar Senha"
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
    );
};

export default UpdatePasswordPage;
