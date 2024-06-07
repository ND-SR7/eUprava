import { jwtDecode } from "jwt-decode";

interface JwtPayload {
  sub: string;
  name: string;
  role: string;
  exp: number;
}

function decodeJwtToken(token: string): JwtPayload {
  try {
    return jwtDecode(token) as JwtPayload;
  } catch(error) {
    console.error("Error while decoding token:", error);
    return {sub: "", name: "", role: "", exp: 0}
  }
  
}

export default decodeJwtToken;
