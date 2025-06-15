/* eslint-disable @next/next/no-img-element */
'use client';
import { redirect, useRouter } from 'next/navigation';
import React, { useContext, useEffect, useState } from 'react';
import { Button } from 'primereact/button';
import { LayoutContext } from '../../../../layout/context/layoutcontext';
import { InputText } from 'primereact/inputtext';
import { classNames } from 'primereact/utils';
import { authApi } from '@/app/services/auth';

const ForgotPasswordPage = () => {
    const [email, setEmail] = useState('');
    const { layoutConfig } = useContext(LayoutContext);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    const router = useRouter();

    const containerClassName = classNames('surface-ground flex align-items-center justify-content-center min-h-screen min-w-screen overflow-hidden', { 'p-input-filled': layoutConfig.inputStyle === 'filled' });

    const handleForgotPassword = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError('');

        const redirect_url = process.env.NEXT_PUBLIC_FORGOT_PASSWORD_REDIRECT_URL;

        try {
            await authApi.forgotPassword({
                email,
                redirect_url: redirect_url
                    ?? 'http://localhost:3000/auth/reset_password'
            });
            router.push('/auth/login');
        } catch (err) {
            router.push('/auth/access');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className={containerClassName}>
            <div className="flex flex-column align-items-center justify-content-center">
                <img src={`/layout/images/logo_juros_facil.png`} alt="Juros Fácil logo" className="mb-5 w-6rem flex-shrink-0" />
                <div
                    style={{
                        borderRadius: '56px',
                        padding: '0.3rem',
                        background: 'rgba(38, 150, 243, 0)'
                    }}
                >
                    <form onSubmit={handleForgotPassword}>
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

                                <div className="flex align-items-center justify-content-between mb-5 gap-5">
                                    <a
                                        className="font-medium no-underline ml-2 text-right cursor-pointer"
                                        style={{ color: 'var(--primary-color)' }}
                                        onClick={() => router.push('/auth/login')}
                                    >
                                        Login
                                    </a>
                                </div>

                                {error && <small className="p-error">{error}</small>}

                                <Button
                                    label="Enviar"
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

export default ForgotPasswordPage;
