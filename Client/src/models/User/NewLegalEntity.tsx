type Role = "USER" | "ADMIN";

type NewLegalEntity = {
  email: string,
  password: string,
  name: string,
  citizenship: string,
  pib: string,
  mb: string,
  role: Role;
  municipality: string,
  locality: string,
  streetName: string,
  streetNumber: number
};

export default NewLegalEntity;
