import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 100,
  duration: '30s',
};

export function setup() {
  const loginPayload = JSON.stringify({
    email: 'fz@gmail.com',
    password: 'password123',
  });

  const loginParams = { headers: { 'Content-Type': 'application/json' } };

  const loginRes = http.post('http://localhost:8081/users/login', loginPayload, loginParams);

  check(loginRes, {
    'login status 200': (r) => r.status === 200,
  });

  const body = loginRes.json();
  const token = body?.results?.token;

  if (!token) {
    throw new Error(`Failed login, Token not found. Response: ${loginRes.body}`);
  }

  return { token }; 
}

export default function (data) {
  const payload = JSON.stringify({
    name: `Book Atomic Habits ${__VU}-${__ITER}`,
    description: 'Atomic Habits by James Clear is a self-help book explaining habit formation, breaking bad habits, and building positive routines.',
    price: 150000,
    stock: 5,
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${data.token}`, 
    },
  };

  const res = http.post('http://localhost:8083/products', payload, params);

  check(res, { 'CreateProduct REST status 201': (r) => r.status === 201 });

  sleep(1);
}
