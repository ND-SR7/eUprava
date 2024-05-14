import { useNavigate } from "react-router-dom";
import Button from "../components/Shared/Button/Button";
import HeadingStyled from "../components/Shared/Heading/Heading.styled";
import { useEffect } from "react";
import { pingCourt } from "../services/PingService";
import toast from "react-hot-toast";

const CourtPage = () => {
  const navigate = useNavigate();
  
  useEffect(() => {
    const token = localStorage.getItem("token") || "";
    if (token === "") navigate("/");
  });

  const ping = () => {
    pingCourt().then(() => {
      toast.success("Good connection with court service");
    }).catch((error) => {
      toast.error("No connection to court service");
      console.error(error);
    });
  };
  
  return (
    <>
      <HeadingStyled>Court</HeadingStyled>
      <br />
      <Button buttonType="button" label="Ping Service" onClick={() => ping()} />
      <br />
      <Button buttonType="button" label="Go Back" onClick={() => window.history.back()}/>
      <br />
    </>
  );
};

export default CourtPage;
