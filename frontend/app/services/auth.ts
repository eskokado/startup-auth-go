import api from "../api/api"

export type RegisterType = {
    name: string
    email: string
    password: string
    passwordConfirmation: string
    imageUrl: string
}

export type ForgotPasswordType = {
    email: string
    redirect_url: string
}

export type UpdatePasswordType = {
    reset_password_token: string
    password: string
    password_confirmation: string
}

export const authApi = {
    // login: async (email: string, password: string) => {
    //     const response = await api.post('/auth/v1/users/sign_in', {
    //         email,
    //         password
    //     });

    //     localStorage.setItem('access-token', response.headers['access-token']);
    //     localStorage.setItem('client', response.headers['client']);
    //     localStorage.setItem('uid', response.headers['uid']);
    //     localStorage.setItem('user-kind', response.headers['user-kind']);
    //     localStorage.setItem('organization', response.headers['organization']);
    //     localStorage.setItem('organization-id', response.headers['organization-id']);
    //     localStorage.setItem('subscription-expires-at', response.headers['subscription-expires-at']);

    //     return response.data;
    // },

    register: async (props: RegisterType) => {
        const response = await api.post('/auth/register', {
            name: props.name,
            email: props.email,
            password: props.password,
            password_confirmation: props.passwordConfirmation,
            image_url: props.imageUrl
        });

        return response.data;
    },

    // forgotPassword: async (props: ForgotPasswordType) => {
    //     const response = await api.post('/auth/v1/users/password', {
    //         email: props.email,
    //         redirect_url: props.redirect_url,
    //     });

    //     return response.data;
    // },

    // updatePassword: async (props: UpdatePasswordType) => {
    //     const response = await api.patch('/auth/v1/users/password', {
    //         reset_password_token: props.reset_password_token,
    //         password: props.password,
    //         password_confirmation: props.password_confirmation,
    //     });

    //     return response.data;
    // },

    // logout: async () => {
    //     await api.delete('/auth/v1/users/sign_out');
    //     localStorage.clear();
    // }
};
