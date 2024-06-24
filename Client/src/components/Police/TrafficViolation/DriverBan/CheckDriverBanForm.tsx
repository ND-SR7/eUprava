import React, { useState } from "react";
import toast from "react-hot-toast";
import HeadingStyled from "../../../Shared/Heading/Heading.styled";
import Button from "../../../Shared/Button/Button";
import { useNavigate } from "react-router-dom";
import { checkDriverBan } from "../../../../services/PoliceService";
import { FormContainer } from "./CheckDriverBanForm.styled";
import { Input } from "../AlcoholLevel/CheckAlcoholLevelForm.styled";

interface CheckDriverBanFormProps {
  closeModal: () => void;
}

const CheckDriverBanForm: React.FC<CheckDriverBanFormProps> = ({ closeModal }) => {
  const [jmbg, setJmbg] = useState("");
  const [location, setLocation] = useState("");
  const navigate = useNavigate();

  const handleSubmit = async (e: any) => {
    e.preventDefault();

    const data = {
      jmbg,
      location,
    };

    try {
      const response = await checkDriverBan(data);
      console.log(response);
      toast.success(response.data || "Driver ban status checked successfully");
      navigate("/home/police");
      closeModal();
    } catch (error: any) {
      console.error(error);
      toast.error("Failed to check driver ban status, check driver JMBG");
    }
  };

  return (
    <FormContainer>
      <HeadingStyled>Check Driver Ban</HeadingStyled>
      <form onSubmit={handleSubmit}>
        <Input
          type="text"
          placeholder="JMBG"
          value={jmbg}
          onChange={(e) => setJmbg(e.target.value)}
          required
        />
        <Input
          type="text"
          placeholder="Location"
          value={location}
          onChange={(e) => setLocation(e.target.value)}
          required
        />
        <Button buttonType="submit" label="Check Driver Ban" />
      </form>
    </FormContainer>
  );
};

export default CheckDriverBanForm;
