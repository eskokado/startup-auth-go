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
}

export const authApi = {
    login: async (email: string, password: string) => {
        const response = await api.post('/auth/login', {
            email,
            password
        });

        return response.data;
    },

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

    forgotPassword: async (props: ForgotPasswordType) => {
        const response = await api.post('/auth/forgot-password', {
            email: props.email,
            redirect_url: props.redirect_url,
        });

        return response.data;
    },

    updatePassword: async (props: UpdatePasswordType) => {
        console.log("erro aqui", props)
        const response = await api.post('/auth/reset-password', {
            reset_password_token: props.reset_password_token,
            password: props.password
        });

        return response.data;
    },

    logout: async () => {
        await api.delete('/auth/logout');
        localStorage.clear();
    }
};
