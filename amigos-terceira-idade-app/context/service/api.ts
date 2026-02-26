import axios from "axios";
import * as SecureStore from "expo-secure-store";
import { Platform } from "react-native";

export const api = axios.create({
  baseURL: "http://localhost:8080/api/v1",
  headers: {
    "Content-Type": "application/json",
  },
});

// 🔐 Interceptor de request (adiciona token automaticamente)
api.interceptors.request.use(async (config) => {
  console.log("joana 123");
  let token: string | null = null;

  if (Platform.OS === "web") {

    token = localStorage.getItem("user_jwt");
  } else {
    token = await SecureStore.getItemAsync("user_jwt");
  }
  console.log(token);

  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  return config;
});

// 🔥 Interceptor de response (logout automático se 401)
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      if (Platform.OS === "web") {
        localStorage.removeItem("user_jwt");
      } else {
        await SecureStore.deleteItemAsync("user_jwt");
      }
    }

    return Promise.reject(error);
  }
);