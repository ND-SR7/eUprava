import { useNavigate } from "react-router-dom";
import Button from "../components/Shared/Button/Button";
import HeadingStyled from "../components/Shared/Heading/Heading.styled";
import { useEffect } from "react";
import { pingStatistics } from "../services/PingService";
import toast from "react-hot-toast";

const StatisticsPage = () => {
  const navigate = useNavigate();
  
  useEffect(() => {
    const token = localStorage.getItem("token") || "";
    if (token === "") navigate("/");
  });

  const ping = () => {
    pingStatistics().then(() => {
      toast.success("Good connection with Institute for Statistics service");
    }).catch((error) => {
      toast.error("No connection to Institute for Statistics service");
      console.error(error);
    });
  };
  
  return (
    <>
      <HeadingStyled>Institute for Statistics</HeadingStyled>
      <br />
      <Button buttonType="button" label="Ping Service" onClick={() => ping()} />
      <br />
      <Button buttonType="button" label="Go Back" onClick={() => window.history.back()}/>
      <br />
    </>
  );
};

export default StatisticsPage;
