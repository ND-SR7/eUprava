import { useNavigate } from "react-router-dom";
import Button from "../components/Shared/Button/Button";
import HeadingStyled from "../components/Shared/Heading/Heading.styled";
import { useEffect } from "react";

const MupPage = () => {
  const navigate = useNavigate();
  
  useEffect(() => {
    const token = localStorage.getItem("token") || "";
    if (token === "") navigate("/");
  });

  return (
    <>
      <HeadingStyled>MUP</HeadingStyled>
      <br />
      <Button buttonType="button" label="Go Back" onClick={() => window.history.back()}/>
      <br />
    </>
  );
};

export default MupPage;
