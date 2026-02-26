import { useRouter } from "expo-router";
import { Image, StyleSheet, Text, TouchableOpacity, View } from "react-native";

export default function MainScreen() {
  const router = useRouter();

  return (
    <View style={styles.container}>
      <Image
        source={require("./assets/logo.png")}
        style={styles.logo}
        resizeMode="contain"
      />

      <Text style={styles.title}>Amigos da Melhor Idade</Text>
      <Text style={styles.subtitle}>
        Conectando pessoas, fortalecendo laços e promovendo bem-estar.
      </Text>

      <TouchableOpacity
        style={styles.primaryButton}
        onPress={() => router.push("/cadastro")}
      >
        <Text style={styles.buttonText}>Criar Conta</Text>
      </TouchableOpacity>

      <TouchableOpacity
        style={styles.secondaryButton}
        onPress={() => router.push("/login")}
      >
        <Text style={styles.secondaryText}>Já tenho conta</Text>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#F4F8FB",
    alignItems: "center",
    justifyContent: "center",
    padding: 24,
  },
  logo: {
    width: 160,
    height: 160,
    marginBottom: 20,
  },
  title: {
    fontSize: 28,
    fontWeight: "bold",
    color: "#1B3A57",
    textAlign: "center",
  },
  subtitle: {
    fontSize: 18,
    color: "#555",
    textAlign: "center",
    marginVertical: 16,
  },
  primaryButton: {
    backgroundColor: "#2E86DE",
    padding: 16,
    borderRadius: 12,
    width: "100%",
    marginTop: 20,
  },
  secondaryButton: {
    marginTop: 12,
  },
  buttonText: {
    color: "#fff",
    fontSize: 20,
    textAlign: "center",
    fontWeight: "bold",
  },
  secondaryText: {
    fontSize: 18,
    color: "#2E86DE",
  },
});