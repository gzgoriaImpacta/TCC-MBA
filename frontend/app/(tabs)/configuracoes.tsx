import { BackendAdapter } from "@/context/adapter/BackendAdapter";
import { BackendService } from "@/context/service/BackedService";
import { useRouter } from "expo-router";
import { useState } from "react";
import { StyleSheet, Switch, Text, TouchableOpacity, View } from "react-native";

export default function Configuracoes() {
  const [darkMode, setDarkMode] = useState(false);
  // const [largeText, setLargeText] = useState(false);
  // const [colorBlindMode, setColorBlindMode] = useState(false);

  const backgroundColor = darkMode ? "#121212" : "#F4F8FB";
  const textColor = darkMode ? "#FFFFFF" : "#000000";

  // const primaryColor = colorBlindMode ? "#0000FF" : "#2E86DE";
  const backendService: BackendAdapter = new BackendService();
  const router = useRouter();
  
  const handleLogout = async () => {
    await backendService.logout();
    router.push("/login");
};

  return (
    <View style={[styles.container, { backgroundColor }]}>
      <Text style={[styles.title, { color: textColor }]}>
        Configurações
      </Text>
      <TouchableOpacity
        onPress={handleLogout}
        style={{
          borderColor: "#3b82f6",
          borderWidth: 1,
          padding: 16,
          borderRadius: 12,
        }}
      >
        <Text
          style={{
            color: "#3b82f6",
            textAlign: "center",
            fontWeight: "bold",
            fontSize: 16,
          }}
        >
          Sair do aplicativo
        </Text>
        </TouchableOpacity>
      {/* <View style={styles.option}>
        <Text style={[styles.label, { color: textColor }]}>
          Modo Escuro
        </Text>
        <Switch value={darkMode} onValueChange={setDarkMode} />
      </View>

      <View style={styles.option}>
        <Text style={[styles.label, { color: textColor }]}>
          Texto Maior
        </Text>
        <Switch value={largeText} onValueChange={setLargeText} />
      </View>

      <View style={styles.option}>
        <Text style={[styles.label, { color: textColor }]}>
          Filtro Daltonismo
        </Text>
        <Switch value={colorBlindMode} onValueChange={setColorBlindMode} />
      </View>

      <Text
        style={{
          fontSize: largeText ? 24 : 18,
          color: primaryColor,
          marginTop: 30,
        }}
      >
        Pré-visualização de Acessibilidade
      </Text> */}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 24,
  },
  title: {
    fontSize: 26,
    fontWeight: "bold",
    marginBottom: 30,
  },
  option: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    marginBottom: 20,
  },
  label: {
    fontSize: 18,
  },
});