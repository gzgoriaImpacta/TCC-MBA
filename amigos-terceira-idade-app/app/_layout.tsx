import { useEffect } from 'react';
import { Stack } from 'expo-router';
import { Provider } from 'react-redux';
import * as SecureStore from 'expo-secure-store';
import { store } from '../store/store';
import { useAppDispatch, useAppSelector } from '../store/hooks';
import { setToken, setLoadingComplete } from '../store/slices/authSlice';
import { ActivityIndicator, View } from 'react-native';
import { MatchProvider } from '../context/MatchContext';
import { ProfileProvider } from '../context/ProfileContext';

// Componente interno para ter acesso aos hooks do Redux
function AppNavigator() {
  const dispatch = useAppDispatch();
  const { isLoading, isAuthenticated } = useAppSelector((state) => state.auth);

  // Efeito para carregar o token quando o app inicia
  useEffect(() => {
    const loadToken = async () => {
      try {
        const storedToken = await SecureStore.getItemAsync('user_jwt');
        if (storedToken) {
          dispatch(setToken(storedToken));
        } else {
          dispatch(setLoadingComplete());
        }
      } catch (error) {
        console.error('Erro ao ler o token:', error);
        dispatch(setLoadingComplete());
      }
    };

    loadToken();
  }, [dispatch]);

  // Tela de carregamento enquanto o token é lido da memória do celular
  if (isLoading) {
    return (
      <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center' }}>
        <ActivityIndicator size="large" />
      </View>
    );
  }

  return (
    <ProfileProvider>
      <MatchProvider>
        <Stack screenOptions={{ headerShown: false }} />
      </MatchProvider>
    </ProfileProvider>
  );
}

export default function RootLayout() {
  return (
    <Provider store={store}>
      <AppNavigator />
    </Provider>
  );
}