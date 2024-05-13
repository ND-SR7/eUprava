import { useNavigate } from "react-router-dom";
import Button from "../components/Shared/Button/Button";
import HeadingStyled from "../components/Shared/Heading/Heading.styled";
import decodeJwtToken from "../services/JwtService";
import { useEffect, useState } from "react";

const HomePage = () => {
  const navigate = useNavigate();
  const [sub, setSub] = useState("");

  useEffect(() => {
    const token = localStorage.getItem("token") || "";
    if (token === "") navigate("/");

    let sub = decodeJwtToken(token).sub;
    sub = sub.substring(0, sub.indexOf("@"))
    setSub(sub);
    // eslint-disable-next-line
  }, []);
  
  return (
    <>
      <HeadingStyled>Welcome to eUprava, {sub}</HeadingStyled>
      <h3>Select the service you want to access</h3>
      <Button
        key="navMUP"
        id="navMUP"
        label="MUP"
        buttonType="button"
        onClick={() => navigate("/home/mup")} />
      <Button
        key="navPolice"
        id="navPolice"
        label="Traffic Police"
        buttonType="button"
        onClick={() => navigate("/home/police")} />
      <Button
        key="navCourt"
        id="navCourt"
        label="Court"
        buttonType="button"
        onClick={() => navigate("/home/court")} />
      <Button
        key="navStatistics"
        id="navStatistics"
        label="Institute for Statistics"
        buttonType="button"
        onClick={() => navigate("/home/statistics")} />
    </>
  );
};

export default HomePage;
