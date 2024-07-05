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
    { label: "Street Number", attrName: "streetNumber", type: "number", min: "0", max: "1000"},
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
    { label: "Street Number", attrName: "streetNumber", type: "number", min: "0", max: "1000"},
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
    // eslint-disable-next-line
    const emailPattern = /^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$/;
    const passwordPattern = /^(?=.*[A-Z])(?=.*[0-9]).{6,}$/;
    const namePattern = /^[A-Za-z]+(?:['\s-][A-Za-z]+)*$/;
    const eighteenYearsAgo = new Date();
    eighteenYearsAgo.setFullYear(eighteenYearsAgo.getFullYear() - 18);
    const jmbgPattern = /^\d{13}$/;
    const pibPattern = /^\d{9}$/;
    const mbPattern = /^\d{8}$/;
    const streetNumberPattern = /^\d+[A-Za-z]*$/;
    const textFieldsPattern = /^[A-Za-z\s-']+$/;

    if (!emailPattern.test(formData["email"])) {
      toast.error("Email format is not valid");
      return;
    }
    
    if (!passwordPattern.test(formData["password"])) {
      toast.error("Password should include at least one uppercase letter, one number and lowercase letters");
      return;
    }
    
    if (!textFieldsPattern.test(formData["citizenship"])) {
      toast.error("Citizenship format is not valid");
      return;
    }

    if (!textFieldsPattern.test(formData["municipality"])) {
      toast.error("Municipality format is not valid");
      return;
    }
    
    if (!textFieldsPattern.test(formData["locality"])) {
      toast.error("Locality format is not valid");
      return;
    }
    
    if (!textFieldsPattern.test(formData["streetName"])) {
      toast.error("Street name format is not valid");
      return;
    }

    if (!streetNumberPattern.test(formData["streetNumber"])) {
      toast.error("Street number format is not valid, must include number and optionaly letters");
      return;
    }

    if (personRegistration) {
      if (!namePattern.test(formData["firstName"])) {
        toast.error("First name format is not valid");
        return;
      }
      
      if (!namePattern.test(formData["lastName"])) {
        toast.error("Last name format is not valid");
        return;
      }

      if (new Date(formData["dob"]) > eighteenYearsAgo) {
        toast.error("You must be 18+ to register");
        return;
      }
  
      if (!jmbgPattern.test(formData["jmbg"])) {
        toast.error("JMBG format is not valid, 13 digits required");
        return;
      }

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
      if (!namePattern.test(formData["name"])) {
        toast.error("Name format is not valid");
        return;
      }

      if (!pibPattern.test(formData["pib"])) {
        toast.error("PIB format is not valid, 9 digits required");
        return;
      }

      if (!mbPattern.test(formData["mb"])) {
        toast.error("MB format is not valid, 8 digits required");
        return;
      }

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
