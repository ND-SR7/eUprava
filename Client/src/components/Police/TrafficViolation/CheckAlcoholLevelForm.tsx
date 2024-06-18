import React, { useState } from "react";
import toast from "react-hot-toast";
import styled from "styled-components";
import HeadingStyled from "../../Shared/Heading/Heading.styled";
import { checkAlcoholLevel } from "../../../services/PoliceService";
import Button from "../../Shared/Button/Button";

const FormContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
`;

const Input = styled.input`
  margin: 10px 0;
  padding: 10px;
  width: 100%;
  max-width: 300px;
`;

const CheckAlcoholLevelForm = () => {
  const [alcoholLevel, setAlcoholLevel] = useState("");
  const [jmbg, setJmbg] = useState("");
  const [location, setLocation] = useState("");

  const handleSubmit = async (e: any) => {
    e.preventDefault();

    const data = {
      alcoholLevel: parseFloat(alcoholLevel),
      jmbg,
      location,
    };

    try {
      const response = await checkAlcoholLevel(data);
      toast.success(response.data);
    } catch (error: any) {
      toast.error("Failed to check alcohol level: " + error.error);
      console.error(error);
    }
  };

  return (
    <FormContainer>
      <HeadingStyled>Check Alcohol Level</HeadingStyled>
      <form onSubmit={handleSubmit}>
        <Input
          type="text"
          placeholder="JMBG"
          value={jmbg}
          onChange={(e) => setJmbg(e.target.value)}
          required
        />
        <Input
          type="number"
          step="0.01"
          placeholder="Alcohol Level"
          value={alcoholLevel}
          onChange={(e) => setAlcoholLevel(e.target.value)}
          required
        />
        <Input
          type="text"
          placeholder="Location"
          value={location}
          onChange={(e) => setLocation(e.target.value)}
          required
        />
        <Button buttonType="submit" label="Check Alcohol Level" />
      </form>
    </FormContainer>
  );
};

export default CheckAlcoholLevelForm;
