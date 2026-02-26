import { UserProfile } from "@/context/adapter/BackendAdapter";
import { BackendService } from "@/context/service/BackedService";
import { useRouter } from "expo-router";
import { useEffect, useState } from "react";
import { Alert, StyleSheet, Text, TextInput, TouchableOpacity, View } from "react-native";

const backendService: BackendAdapter = new BackendService();

export default function Perfil() {
  const router = useRouter();

  const [nome, setNome] = useState("");
  const [email, setEmail] = useState("");
  const [usuario, setUsuario] = useState("");
  const [idade, setIdade] = useState("");
  const [tipo, setTipo] = useState("");
  const [interesses, setInteresses] = useState("");

  const [profile, setProfile] = useState<UserProfile | null>(null);

  useEffect(() => {
    async function loadProfile() {
        try {
          const data = await backendService.getUserProfile();
          setProfile(data);
        } catch (error) {
          console.error(error);
        }
        console.log("joana")
      }
    loadProfile();
  }, []);


  const salvarPerfil = () => {
    if (!nome || !email) {
      Alert.alert("Preencha os campos obrigatórios");
      return;
    }

    Alert.alert("Perfil salvo com sucesso!");
    router.replace("/(tabs)");
  };

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Meu Perfil</Text>

      <TextInput
        placeholder="Nome completo"
        style={styles.input}
        value={nome}
        onChangeText={setNome}
      />

      <TextInput
        placeholder="Email"
        style={styles.input}
        value={email}
        onChangeText={setEmail}
      />

      <TextInput
        placeholder="Nome de usuário"
        style={styles.input}
        value={usuario}
        onChangeText={setUsuario}
      />

      <TextInput
        placeholder="Idade"
        keyboardType="numeric"
        style={styles.input}
        value={idade}
        onChangeText={setIdade}
      />

      <TextInput
        placeholder="Tipo (Idoso ou Voluntário)"
        style={styles.input}
        value={tipo}
        onChangeText={setTipo}
      />

      <TextInput
        placeholder="Interesses (ex: caminhada, leitura)"
        style={styles.input}
        value={interesses}
        onChangeText={setInteresses}
      />

      <TouchableOpacity style={styles.button} onPress={salvarPerfil}>
        <Text style={styles.buttonText}>Salvar Perfil</Text>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 24,
    backgroundColor: "#F4F8FB",
  },
  title: {
    fontSize: 26,
    fontWeight: "bold",
    marginBottom: 20,
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