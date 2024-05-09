import toast from "react-hot-toast";
import Form from "../components/Shared/Form/Form";
import { useNavigate } from "react-router-dom";
import Button from "../components/Shared/Button/Button";
import UserToken from "../models/User/UserToken";
import { login } from "../services/SSOService";

const LoginPage = () => {
  const navigate = useNavigate();
  
  const formFields = [
    { label: "Email", attrName: "email", type: "text"},
    { label: "Password", attrName: "password", type: "password"}
  ];

  function loginUser(formData: any): void {
    login({
      email: formData["email"],
      password: formData["password"]
    }).then((userToken: UserToken) => {
      localStorage.setItem("token", userToken.token);
      toast.success("Successfully logged in");
      navigate("/todo");
    }).catch(() => {
      toast.error("Failed to log in.\nCheck credentials or activate your account");
    });
  }

  return (
    <>
      <Form formFields={formFields} heading="Login" onSubmit={loginUser} />
      <br />
      <Button key="btnRegister" id="btnRegister" label="Don't have an account? Click here to register" buttonType="button" onClick={() => console.log("ToDo")} />
      <br />
      <Button key="btnRecovery" id="btnRecovery" label="Forgot your password? Click here to recover it" buttonType="button" onClick={() => console.log("ToDo")} />
    </>
  );
};

export default LoginPage;
