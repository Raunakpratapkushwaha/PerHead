import React, { useState } from 'react';
import { View, Text, TextInput, TouchableOpacity, StyleSheet, ActivityIndicator } from 'react-native';
import { useAuthStore } from '../../src/store/authStore';

export default function LoginScreen() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  
  // Bring in our global state and login function
  const { login, isLoading, error, user } = useAuthStore();

  const handleLogin = () => {
    login(email, password);
  };

  // If we successfully logged in, show a success message!
  if (user) {
    return (
      <View style={styles.container}>
        <Text style={styles.title}>Welcome to Batwara, {user.name}! 🚀</Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Batwara</Text>
      <Text style={styles.subtitle}>Settle up with friends.</Text>

      {error && <Text style={styles.errorText}>{error}</Text>}

      <TextInput
        style={styles.input}
        placeholder="Email"
        value={email}
        onChangeText={setEmail}
        autoCapitalize="none"
        keyboardType="email-address"
      />
      
      <TextInput
        style={styles.input}
        placeholder="Password"
        value={password}
        onChangeText={setPassword}
        secureTextEntry
      />

      <TouchableOpacity style={styles.button} onPress={handleLogin} disabled={isLoading}>
        {isLoading ? (
          <ActivityIndicator color="#fff" />
        ) : (
          <Text style={styles.buttonText}>Log In</Text>
        )}
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, justifyContent: 'center', padding: 20, backgroundColor: '#f8f9fa' },
  title: { fontSize: 32, fontWeight: 'bold', textAlign: 'center', color: '#1a1a1a', marginBottom: 5 },
  subtitle: { fontSize: 16, textAlign: 'center', color: '#6c757d', marginBottom: 30 },
  input: { backgroundColor: '#fff', padding: 15, borderRadius: 10, marginBottom: 15, borderWidth: 1, borderColor: '#dee2e6' },
  button: { backgroundColor: '#0d6efd', padding: 15, borderRadius: 10, alignItems: 'center', marginTop: 10 },
  buttonText: { color: '#fff', fontSize: 16, fontWeight: 'bold' },
  errorText: { color: '#dc3545', textAlign: 'center', marginBottom: 15, fontWeight: '500' }
});