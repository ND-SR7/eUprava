import toast from "react-hot-toast";
import Form from "../components/Shared/Form/Form";
import { registerPerson, registerLegalEntity } from "../services/SSOService";
import { useNavigate } from "react-router-dom";
import Button from "../components/Shared/Button/Button";
import { useState } from "react";

const RegisterPage = () => {
  const navigate = useNavigate();
  const [personRegistration, setPersonRegistration] = useState(true);
  
  const formFieldsPerson = [
    { label: "First Name", attrName: "firstName", type: "text"},
    { label: "Last Name", attrName: "lastName", type: "text"},
    { label: "Sex", attrName: "sex", type: "radio", options: ["MALE", "FEMALE"]},
    { label: "Citizenship", attrName: "citizenship", type: "text"},
    { label: "Date of Birth", attrName: "dob", type: "date"},
    { label: "JMBG", attrName: "jmbg", type: "text"},
    { label: "Municipality", attrName: "municipality", type: "text"},
    { label: "Locality", attrName: "locality", type: "text"},
    { label: "Street Name", attrName: "streetName", type: "text"},
    { label: "Street Number", attrName: "streetNumber", type: "number"},
    { label: "Email", attrName: "email", type: "email"},
    { label: "Password", attrName: "password", type: "password"},
    { label: "Repeat Password", attrName: "passwordRepeat", type: "password"},
  ];

  const formFieldsLegalEntity = [
    { label: "Name", attrName: "name", type: "text"},
    { label: "Citizenship", attrName: "citizenship", type: "text"},
    { label: "PIB", attrName: "pib", type: "text"},
    { label: "MB", attrName: "mb", type: "text"},
    { label: "Municipality", attrName: "municipality", type: "text"},
    { label: "Locality", attrName: "locality", type: "text"},
    { label: "Street Name", attrName: "streetName", type: "text"},
    { label: "Street Number", attrName: "streetNumber", type: "number"},
    { label: "Email", attrName: "email", type: "email"},
    { label: "Password", attrName: "password", type: "password"},
    { label: "Repeat Password", attrName: "passwordRepeat", type: "password"},
  ];

  function setPassword(password: string, repeatedPassword: string): string {
    if (password === repeatedPassword) {
      return password;
    }
    toast.error("Passwords do not match");
    return "";
  }

  function registerUser(formData: any) {
    // TODO: Validation
    console.log(formData);

    if (personRegistration) {
      registerPerson({
        email: formData["email"],
        password: setPassword(formData["password"], formData["passwordRepeat"]),
        firstName: formData["firstName"],
        lastName: formData["lastName"],
        sex: formData["sex"],
        citizenship: formData["citizenship"],
        dob: formData["dob"],
        jmbg: formData["jmbg"],
        role: "USER",
        municipality: formData["municipality"],
        locality: formData["locality"],
        streetName: formData["streetName"],
        streetNumber: Number(formData["streetNumber"])
      }).then(() => {
        toast.success("Successfully registered");
        navigate("/");
      }).catch(() => {
        toast.error("Failed to register user");
      });
    } else {
      registerLegalEntity({
        email: formData["email"],
        password: setPassword(formData["password"], formData["passwordRepeat"]),
        name: formData["name"],
        citizenship: formData["citizenship"],
        pib: formData["pib"],
        mb: formData["mb"],
        role: "USER",
        municipality: formData["municipality"],
        locality: formData["locality"],
        streetName: formData["streetName"],
        streetNumber: Number(formData["streetNumber"])
      }).then(() => {
        toast.success("Successfully registered");
        navigate("/");
      }).catch(() => {
        toast.error("Failed to register user");
      });
    }
  }

  return (
    <>
      <Button
        key="btnSwitchRegMode"
        id="btnSwitchRegMode"
        label={personRegistration ? "Switch to legal entity form" : "Switch to person form"}
        buttonType="button"
        onClick={() => setPersonRegistration(!personRegistration)} />
      <Form
        formFields={personRegistration ? formFieldsPerson : formFieldsLegalEntity}
        heading="Register"
        onSubmit={registerUser} />
      <br />
      <Button
        key="btnRegister"
        id="btnRegister"
        label="Already have an account? Click here to login"
        buttonType="button"
        onClick={() => navigate("/")} />
    </>
  );
};

export default RegisterPage;
