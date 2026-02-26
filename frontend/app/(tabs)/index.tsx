import { useRef, useState } from "react";
import {
  Animated,
  Dimensions,
  Modal,
  PanResponder,
  StyleSheet,
  Text,
  TouchableOpacity,
  View,
} from "react-native";

const { width } = Dimensions.get("window");

type Perfil = {
  id: string;
  nome: string;
  idade: number;
  interesse: string;
};

export default function Home() {
  const [perfis, setPerfis] = useState<Perfil[]>([
    { id: "1", nome: "Maria", idade: 72, interesse: "Gosta de caminhadas" },
    { id: "2", nome: "João", idade: 68, interesse: "Ama jogar dominó" },
    { id: "3", nome: "Ana", idade: 74, interesse: "Adora conversar e ler" },
  ]);

  const [matches, setMatches] = useState<Perfil[]>([]);
  const [modalVisible, setModalVisible] = useState(false);
  const [perfilMatch, setPerfilMatch] = useState<Perfil | null>(null);

  const position = useRef(new Animated.ValueXY()).current;

  const rotate = position.x.interpolate({
    inputRange: [-width / 2, 0, width / 2],
    outputRange: ["-15deg", "0deg", "15deg"],
    extrapolate: "clamp",
  });

  const likeOpacity = position.x.interpolate({
    inputRange: [0, width / 4],
    outputRange: [0, 1],
    extrapolate: "clamp",
  });

  const nopeOpacity = position.x.interpolate({
    inputRange: [-width / 4, 0],
    outputRange: [1, 0],
    extrapolate: "clamp",
  });

  const removerPerfil = (curtiu: boolean) => {
    const atual = perfis[0];

    if (curtiu) {
      setMatches((prev) => [...prev, atual]);
      setPerfilMatch(atual);
      setModalVisible(true);
    }

    setPerfis((prev) => prev.slice(1));
    position.setValue({ x: 0, y: 0 });
  };

  const panResponder = PanResponder.create({
    onStartShouldSetPanResponder: () => true,

    onPanResponderMove: (_, gesture) => {
      position.setValue({ x: gesture.dx, y: gesture.dy });
    },

    onPanResponderRelease: (_, gesture) => {
      if (gesture.dx > 120) {
        Animated.timing(position, {
          toValue: { x: width + 100, y: gesture.dy },
          duration: 250,
          useNativeDriver: false,
        }).start(() => removerPerfil(true));
      } else if (gesture.dx < -120) {
        Animated.timing(position, {
          toValue: { x: -width - 100, y: gesture.dy },
          duration: 250,
          useNativeDriver: false,
        }).start(() => removerPerfil(false));
      } else {
        Animated.spring(position, {
          toValue: { x: 0, y: 0 },
          useNativeDriver: false,
        }).start();
      }
    },
  });

  if (perfis.length === 0) {
    return (
      <View style={styles.container}>
        <Text style={styles.empty}>Sem novos perfis por hoje 💙</Text>
        <Text style={styles.sub}>
          Matches realizados: {matches.length}
        </Text>
      </View>
    );
  }

  const perfilAtual = perfis[0];

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Encontrar Companhia</Text>

      <Animated.View
        {...panResponder.panHandlers}
        style={[
          styles.card,
          {
            transform: [
              { translateX: position.x },
              { translateY: position.y },
              { rotate },
            ],
          },
        ]}
      >
        <Animated.Text style={[styles.like, { opacity: likeOpacity }]}>
          CURTIDO 💚
        </Animated.Text>

        <Animated.Text style={[styles.nope, { opacity: nopeOpacity }]}>
          NÃO ❌
        </Animated.Text>

        <Text style={styles.nome}>
          {perfilAtual.nome}, {perfilAtual.idade}
        </Text>
        <Text style={styles.interesse}>{perfilAtual.interesse}</Text>
      </Animated.View>

      <Modal visible={modalVisible} transparent animationType="fade">
        <View style={styles.modalContainer}>
          <View style={styles.modalBox}>
            <Text style={styles.matchTitle}>🎉 É um Match!</Text>
            <Text style={styles.matchText}>
              Você e {perfilMatch?.nome} demonstraram interesse.
            </Text>

            <TouchableOpacity
              style={styles.button}
              onPress={() => setModalVisible(false)}
            >
              <Text style={styles.buttonText}>Continuar</Text>
            </TouchableOpacity>
          </View>
        </View>
      </Modal>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#F4F8FB",
    alignItems: "center",
    paddingTop: 40,
  },
  title: {
    fontSize: 24,
    fontWeight: "bold",
    marginBottom: 30,
  },
  card: {
    width: width * 0.9,
    height: 400,
    backgroundColor: "#fff",
    borderRadius: 20,
    justifyContent: "center",
    alignItems: "center",
    padding: 20,
    elevation: 6,
  },
  nome: {
    fontSize: 26,
    fontWeight: "bold",
  },
  interesse: {
    fontSize: 18,
    marginTop: 10,
    textAlign: "center",
  },
  like: {
    position: "absolute",
    top: 40,
    left: 20,
    fontSize: 22,
    fontWeight: "bold",
    color: "green",
  },
  nope: {
    position: "absolute",
    top: 40,
    right: 20,
    fontSize: 22,
    fontWeight: "bold",
    color: "red",
  },
  empty: {
    fontSize: 20,
    marginTop: 200,
  },
  sub: {
    marginTop: 10,
    fontSize: 16,
  },
  modalContainer: {
    flex: 1,
    backgroundColor: "rgba(0,0,0,0.5)",
    justifyContent: "center",
    alignItems: "center",
  },
  modalBox: {
    width: "80%",
    backgroundColor: "#fff",
    padding: 30,
    borderRadius: 20,
    alignItems: "center",
  },
  matchTitle: {
    fontSize: 26,
    fontWeight: "bold",
    marginBottom: 10,
  },
  matchText: {
    fontSize: 18,
    textAlign: "center",
    marginBottom: 20,
  },
  button: {
    backgroundColor: "#4CAF50",
    paddingVertical: 12,
    paddingHorizontal: 30,
    borderRadius: 10,
  },
  buttonText: {
    color: "#fff",
    fontSize: 16,
  },
});