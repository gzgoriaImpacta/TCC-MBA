import { useRouter } from "expo-router";
import { Platform, StyleSheet, Text, TextInput, TouchableOpacity, View } from "react-native";
import { BackendAdapter, LoginRequest } from "../context/adapter/BackendAdapter";
import { BackendService } from "@/context/service/BackedService";
import { useState } from "react";
import { useAppDispatch } from "@/store/hooks";
import * as SecureStore from 'expo-secure-store';
import { setToken } from "@/store/slices/authSlice";

const backendService: BackendAdapter = new BackendService();

export default function Login() {
  const dispatch = useAppDispatch();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const router = useRouter();

  const handleLogin = async () => {
    // Lógica de autenticação aqui
    const loginRequest: LoginRequest = {
      email,
      password,
    };

    try {
      const response = await backendService.login(loginRequest);
      if (response.data.access_token) {
        try {

          if (Platform.OS === 'web') {
            localStorage.setItem("refresh_token", response.data.refresh_token);
            localStorage.setItem('user_jwt', response.data.access_token);
          } else {
            await SecureStore.setItemAsync('user_jwt', response.data.access_token);
            await SecureStore.setItemAsync("refresh_token", response.data.refresh_token);
          }
        } catch (secureStoreError) {
          console.warn('SecureStore error, falling back to AsyncStorage:', secureStoreError);
          // Fallback em caso de erro com SecureStore
        }
        dispatch(setToken(response.data.access_token));
        router.replace("/(tabs)");
      } else {
        alert("Login falhou: Token não recebido");
      }
    } catch (error) {
      console.error('Login error:', error);
      alert("Erro ao fazer login");
    }
  }

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Entrar</Text>

      <TextInput placeholder="Email" style={styles.input} onChange={(e) => setEmail(e.nativeEvent.text)} />
      <TextInput placeholder="Senha" secureTextEntry style={styles.input} onChange={(e) => setPassword(e.nativeEvent.text)} />

      <TouchableOpacity
        style={styles.button}
        onPress={() => handleLogin()}
      >
        <Text style={styles.buttonText}>Entrar</Text>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 24,
    backgroundColor: "#F4F8FB",
    justifyContent: "center",
  },
  title: {
    fontSize: 26,
    fontWeight: "bold",
    marginBottom: 20,
    textAlign: "center",
  },
  input: {
    backgroundColor: "#fff",
    padding: 16,
    borderRadius: 10,
    marginBottom: 14,
    fontSize: 18,
  },
  button: {
    backgroundColor: "#2E86DE",
    padding: 16,
    borderRadius: 12,
    marginTop: 10,
  },
  buttonText: {
    color: "#fff",
    fontSize: 20,
    textAlign: "center",
    fontWeight: "bold",
  },
});