import { useNavigate } from "react-router-dom";
import Button from "../components/Shared/Button/Button";
import HeadingStyled from "../components/Shared/Heading/Heading.styled";
import decodeJwtToken from "../services/JwtService";
import { useEffect, useState } from "react";

const HomePage = () => {
  const navigate = useNavigate();
  const [name, setName] = useState("");

  useEffect(() => {
    const token = localStorage.getItem("token") || "";
    if (token === "") navigate("/");

    let tokenName = decodeJwtToken(token).name;
    setName(tokenName);
    // eslint-disable-next-line
  }, []);

  const logout = () => {
    localStorage.removeItem("token");
    navigate("/");
  };

  return (
    <>
      <HeadingStyled>Welcome to eUprava, {name}</HeadingStyled>
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
        onClick={() => navigate("/home/police")}
      />
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
      <br />
      <br />
      <Button
        key="btnLogout"
        id="btnLogout"
        label="Logout"
        buttonType="button"
        onClick={() => logout()} />
    </>
  );
};

export default HomePage;
