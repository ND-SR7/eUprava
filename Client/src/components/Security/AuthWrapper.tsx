import decodeJwtToken from "../../services/JwtService";

export enum UserRole {
  User = 'USER',
  Admin = 'ADMIN',
}

interface AuthWrapperProps {
  allowedRoles: UserRole[];
  children: React.ReactNode;
}

const AuthWrapper = ({ allowedRoles, children }: AuthWrapperProps) => {
  const token = localStorage.getItem('token');
  
  if (!token) return null;

  const decodedToken = decodeJwtToken(token);
  const userRole = decodedToken.role as UserRole;

  if (!allowedRoles.includes(userRole)) return null;

  return <>{children}</>;
};

export default AuthWrapper;
