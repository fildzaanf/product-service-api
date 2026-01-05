import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';

const userClient = new grpc.Client();
const productClient = new grpc.Client();

userClient.load(['proto'], 'user.proto');       
productClient.load(['proto'], 'product.proto'); 

export const options = {
  vus: 100,         
  duration: '30s',  
};


export function setup() {
  userClient.connect('localhost:8080', { plaintext: true }); 

  const loginRes = userClient.invoke('user.UserCommandService/LoginUser', {
    email: 'fz@gmail.com',
    password: 'password123',
  });

  check(loginRes, {
    'login status OK': (r) => r && r.status === grpc.StatusOK,
    'token returned': (r) => r && r.message && r.message.token,
  });

  userClient.close();

  if (!loginRes.message || !loginRes.message.token) {
    throw new Error('Failed login, Token not found');
  }

  return { token: loginRes.message.token };
}


export default function (data) {
  productClient.connect('localhost:8082', { plaintext: true }); 

  const payload = {
    name: `Book Atomic Habits ${__VU}-${__ITER}`,
    description:
      'Atomic Habits by James Clear is a self-help book explaining habit formation, breaking bad habits, and building positive routines.',
    price: '150000',
    stock: 5,
  };

  const res = productClient.invoke(
    'product.ProductCommandService/CreateProduct',
    payload,
    {
      metadata: {
        authorization: `Bearer ${data.token}`, 
      },
    }
  );

  check(res, {
    'CreateProduct OK': (r) => r && r.status === grpc.StatusOK,
  });

  productClient.close();
  sleep(1); 
}
