import axios from 'axios';

export const baseURL = process.env.NEXT_PUBLIC_API_URL

const api = axios.create({ baseURL, withCredentials: true });

// api.interceptors.request.use((config) => {
//     const token = localStorage.getItem('access-token');
//     if (token) {
//         config.headers['access-token'] = token;
//         config.headers['client'] = localStorage.getItem('client');
//         config.headers['uid'] = localStorage.getItem('uid');
//     }
//     return config;
// });

// api.interceptors.response.use(
//     (response) => {
//         // Atualiza os tokens no localStorage se novos headers estiverem presentes
//         const newAccessToken = response.headers['access-token'];
//         const newClient = response.headers['client'];
//         const newUid = response.headers['uid'];
//         const userKind = response.headers['user-kind'];
//         const organization = response.headers['organization'];
//         const organizationId = response.headers['organization-id'];
//         const subscriptionExpiresAt = response.headers['subscription-expires-at'];

//         if (newAccessToken && newClient && newUid) {
//             localStorage.setItem('access-token', newAccessToken);
//             localStorage.setItem('client', newClient);
//             localStorage.setItem('uid', newUid);
//             localStorage.setItem('user-kind', userKind);
//             localStorage.setItem('organization', organization);
//             localStorage.setItem('organization-id', organizationId);
//             localStorage.setItem('subscription-expires-at', subscriptionExpiresAt)
//         }
//         // if (response.data?.data?.subscription_expires_at) {
//         //     localStorage.setItem('subscription_expires_at', response.data.data.subscription_expires_at);
//         // }

//         return response;
//     },
//     (error) => {
//         if (error.response) {
//             const subscriptionExpiresAt = error.response.headers['subscription-expires-at'];
//             if (subscriptionExpiresAt) {
//                 localStorage.setItem('subscription_expires_at', subscriptionExpiresAt);
//                 // Verifica se a assinatura expirou
//                 if (new Date(subscriptionExpiresAt) < new Date()) {
//                     localStorage.removeItem('access-token');
//                     localStorage.removeItem('client');
//                     localStorage.removeItem('uid');
//                     localStorage.removeItem('subscription-expires-at');
//                     localStorage.removeItem('user-kind');
//                     localStorage.removeItem('organization');
//                     localStorage.removeItem('organization-id');
//                     window.location.href = '/landing';
//                 }
//             }

//             if (error.response?.status === 401) {
//                 // Remove tokens do localStorage
//                 localStorage.removeItem('access-token');
//                 localStorage.removeItem('client');
//                 localStorage.removeItem('uid');
//                 localStorage.removeItem('subscription-expires-at');
//                 localStorage.removeItem('user-kind');
//                 localStorage.removeItem('organization');
//                 localStorage.removeItem('organization-id');
//                 window.location.href = '/landing';
//                 return Promise.reject({ code: 401, messages: ['Não autorizado'] });
//             }


//             const errorData = error.response.data;
//             let errorMessages: string[] = [];

//             if (errorData.errors && Array.isArray(errorData.errors)) {
//                 errorMessages = errorData.errors; // Array
//             } else if (errorData.error) {
//                 errorMessages = [errorData.error]; // Converte string em array
//             } else {
//                 errorMessages = ['Erro desconhecido'];
//             }

//             const errorCode = error.response.status || 500;
//             return Promise.reject({ code: errorCode, messages: errorMessages });
//         }

//         const errorMessage = error.response?.data?.message || 'Erro desconhecido';
//         const errorCode = error.response?.status || 500;
//         return Promise.reject({ code: errorCode, messages: [errorMessage] });
//     }
// );

export default api;
