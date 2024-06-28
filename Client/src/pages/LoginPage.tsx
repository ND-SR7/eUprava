import toast from "react-hot-toast";
import Form from "../components/Shared/Form/Form";
import { useNavigate } from "react-router-dom";
import Button from "../components/Shared/Button/Button";
import UserToken from "../models/User/UserToken";
import { login, sendRecoveryEmail } from "../services/SSOService";
import { useEffect, useState } from "react";
import Modal from "../components/Shared/Modal/Modal";

const LoginPage = () => {
  const navigate = useNavigate();
  const [recoveryModal, setRecoveryModal] = useState(false);

  const recoveryForm = (
    <Form
      heading=""
      formFields={[{ label: "Email", attrName: "email", type: "text"}]}
      onSubmit={getRecoveryEmail} />
  );
  
  const formFields = [
    { label: "Email", attrName: "email", type: "text"},
    { label: "Password", attrName: "password", type: "password"}
  ];

  useEffect(() => {
    if (localStorage.getItem("token") !== null) navigate("/home");
  });

  function loginUser(formData: any): void {
    // eslint-disable-next-line
    const emailPattern = /^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$/;
    const passwordPattern = /^(?=.*[A-Z])(?=.*[0-9]).{6,}$/;

    if (!emailPattern.test(formData["email"])) {
      toast.error("Email format is not valid");
      return;
    }
    
    if (!passwordPattern.test(formData["password"])) {
      toast.error("Password should include at least one uppercase letter, one number and lowercase letters");
      return;
    }

    login({
      email: formData["email"],
      password: formData["password"]
    }).then((userToken: UserToken) => {
      localStorage.setItem("token", userToken.token);
      toast.success("Successfully logged in");
      navigate("/home");
    }).catch(() => {
      toast.error("Failed to log in.\nCheck credentials or activate your account");
    });
  }

  function getRecoveryEmail(formData: any): void {
    sendRecoveryEmail(formData["email"]).then(() => {
      toast.success("Recovery email sent. Check you inbox");
      setRecoveryModal(false);
    }).catch((error) => {
      console.error(error);
      toast.error("Failed to send recovery email")
    })
  }

  return (
    <>
      <Form formFields={formFields} heading="Login" onSubmit={loginUser} />
      <br />
      <Button key="btnRegister" id="btnRegister" label="Don't have an account? Click here to register" buttonType="button" onClick={() => navigate("/register")} />
      <br />
      <Button key="btnRecovery" id="btnRecovery" label="Forgot your password? Click here to recover it" buttonType="button" onClick={() => setRecoveryModal(true)} />

      <Modal heading="Enter email for password recovery" content={recoveryForm} isVisible={recoveryModal} onClose={() => setRecoveryModal(false)} />
    </>
  );
};

export default LoginPage;
