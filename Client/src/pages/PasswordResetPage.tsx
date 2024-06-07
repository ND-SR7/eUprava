import toast from "react-hot-toast";
import Form from "../components/Shared/Form/Form";
import HeadingStyled from "../components/Shared/Heading/Heading.styled";
import { resetPassword } from "../services/SSOService";
import { useNavigate } from "react-router-dom";

const PasswordResetPage = () => {
  const navigate = useNavigate();

  const formFields = [
    { label: "Password reset code", attrName: "passwordResetCode", type: "text"},
    { label: "New password", attrName: "newPassword", type: "password"},
    { label: "Confirm new password", attrName: "confirmNewPassword", type: "password"},
  ];

  const sendPasswordReset = (formData: any): void => {
    if (formData["newPassword"] !== formData["confirmNewPassword"]) {
      toast.error("Passwords do not match");
      return;
    }
    
    resetPassword({
      passwordResetCode: formData["passwordResetCode"],
      newPassword: formData["newPassword"]
    }).then(() => {
      toast.success("Password successfully reset");
      navigate("/");
    }).catch((error) => {
      toast.error("Failed to reset password");
      console.error(error);
    });
  };

  return (
    <>
      <HeadingStyled>Password Reset</HeadingStyled>
      <Form
        heading="Enter reset code and new password"
        formFields={formFields}
        onSubmit={sendPasswordReset} />
    </>
  );
};

export default PasswordResetPage;
