type Role = "USER" | "ADMIN";
type Sex = "MALE" | "FEMALE";

type NewPerson = {
  email: string,
  password: string,
  firstName: string,
  lastName: string,
  sex: Sex,
  citizenship: string,
  dob: string,
  jmbg: string,
  role: Role;
  municipality: string,
  locality: string,
  streetName: string,
  streetNumber: number
};

export default NewPerson;
