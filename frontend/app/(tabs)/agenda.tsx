// import { useState } from "react";
// import {
//     Animated,
//     FlatList,
//     StyleSheet,
//     Text,
//     TouchableOpacity,
//     View,
// } from "react-native";

// type Atividade = {
//   id: string;
//   titulo: string;
//   data: string;
//   confirmada: boolean;
//   lida: boolean;
// };

// function ItemAgenda({
//   item,
//   onConfirmar,
//   onLer,
// }: {
//   item: Atividade;
//   onConfirmar: () => void;
//   onLer: () => void;
// }) {
//   const scale = new Animated.Value(1);

//   const animarConfirmacao = () => {
//     Animated.sequence([
//       Animated.timing(scale, {
//         toValue: 1.1,
//         duration: 150,
//         useNativeDriver: true,
//       }),
//       Animated.timing(scale, {
//         toValue: 1,
//         duration: 150,
//         useNativeDriver: true,
//       }),
//     ]).start(onConfirmar);
//   };

//   return (
//     <Animated.View
//       style={[
//         styles.card,
//         {
//           transform: [{ scale }],
//           backgroundColor: item.lida ? "#E8F5E9" : "#fff",
//         },
//       ]}
//     >
//       <Text style={styles.titulo}>{item.titulo}</Text>
//       <Text style={styles.data}>{item.data}</Text>

//       <View style={styles.buttons}>
//         {!item.confirmada && (
//           <TouchableOpacity
//             style={styles.confirmar}
//             onPress={animarConfirmacao}
//           >
//             <Text style={styles.btnText}>Confirmar</Text>
//           </TouchableOpacity>
//         )}

//         {!item.lida && (
//           <TouchableOpacity style={styles.ler} onPress={onLer}>
//             <Text style={styles.btnText}>Marcar como lida</Text>
//           </TouchableOpacity>
//         )}
//       </View>
//     </Animated.View>
//   );
// }

// export default function Agenda() {
//   const [atividades, setAtividades] = useState<Atividade[]>([
//     {
//       id: "1",
//       titulo: "Caminhada no Parque",
//       data: "10/02 - 08:00",
//       confirmada: false,
//       lida: false,
//     },
//     {
//       id: "2",
//       titulo: "Jogo de Dominó",
//       data: "11/02 - 14:00",
//       confirmada: false,
//       lida: false,
//     },
//   ]);

//   const confirmar = (id: string) => {
//     setAtividades((prev) =>
//       prev.map((item) =>
//         item.id === id ? { ...item, confirmada: true } : item
//       )
//     );
//   };

//   const marcarComoLida = (id: string) => {
//     setAtividades((prev) =>
//       prev.map((item) =>
//         item.id === id ? { ...item, lida: true } : item
//       )
//     );
//   };

//   return (
//     <View style={styles.container}>
//       <Text style={styles.header}>Minha Agenda</Text>

//       <FlatList
//         data={atividades}
//         keyExtractor={(item) => item.id}
//         renderItem={({ item }) => (
//           <ItemAgenda
//             item={item}
//             onConfirmar={() => confirmar(item.id)}
//             onLer={() => marcarComoLida(item.id)}
//           />
//         )}
//       />
//     </View>
//   );
// }

// const styles = StyleSheet.create({
//   container: {
//     flex: 1,
//     padding: 20,
//     backgroundColor: "#F4F8FB",
//   },
//   header: {
//     fontSize: 24,
//     fontWeight: "bold",
//     marginBottom: 20,
//   },
//   card: {
//     padding: 20,
//     borderRadius: 15,
//     marginBottom: 15,
//     elevation: 4,
//   },
//   titulo: {
//     fontSize: 20,
//     fontWeight: "bold",
//   },
//   data: {
//     fontSize: 16,
//     marginVertical: 5,
//   },
//   buttons: {
//     marginTop: 10,
//     flexDirection: "row",
//     gap: 10,
//   },
//   confirmar: {
//     backgroundColor: "#4CAF50",
//     padding: 10,
//     borderRadius: 8,
//   },
//   ler: {
//     backgroundColor: "#2196F3",
//     padding: 10,
//     borderRadius: 8,
//   },
//   btnText: {
//     color: "#fff",
//     fontWeight: "bold",
//   },
// });