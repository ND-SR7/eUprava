import React, { useState } from 'react';
import toast from 'react-hot-toast';
import { FormContainer, Input, Select, StyledHeading, SubmitButton } from './CheckAllForm.styled';
import { useNavigate } from 'react-router-dom';
import { checkAll } from '../../../services/PoliceService';

interface CheckDriverFormProps {
  closeModal: () => void;
}

const CheckDriverForm: React.FC<CheckDriverFormProps> = ({ closeModal }) => {
  const [jmbg, setJmbg] = useState('');
  const [alcoholLevel, setAlcoholLevel] = useState('');
  const [tire, setTire] = useState('SUMMER');
  const [platesNumber, setPlatesNumber] = useState('');
  const [location, setLocation] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const data = {
      jmbg,
      alcoholLevel: parseFloat(alcoholLevel),
      tire,
      platesNumber,
      location,
    };

    console.log("Submitting data:", data);

    try {
      const response = await checkAll(data);
      console.log("API response:", response);
  
      const message = response?.message;
      if (message) {
        toast.success(message);
      } else {
        toast.success('Driver checked successfully');
      }
      
      navigate('/home/police');
      closeModal();
    } catch (error: any) {
      console.error("Error response:", error);
  
      const errorMessage = error?.response?.data?.message;
      if (errorMessage) {
        toast.error(errorMessage);
      } else {
        toast.error('Failed to check driver, check all input fields');
      }
    }
  };

  return (
    <FormContainer>
      <StyledHeading>Check Driver</StyledHeading>
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
        <Select
          value={tire}
          onChange={(e) => setTire(e.target.value)}
          required
        >
          <option value="SUMMER">SUMMER</option>
          <option value="WINTER">WINTER</option>
        </Select>
        <Input
          type="text"
          placeholder="Plates Number"
          value={platesNumber}
          onChange={(e) => setPlatesNumber(e.target.value)}
          required
        />
        <Input
          type="text"
          placeholder="Location"
          value={location}
          onChange={(e) => setLocation(e.target.value)}
          required
        />
        <SubmitButton buttonType="submit" label="Check Driver" />
      </form>
    </FormContainer>
  );
};

export default CheckDriverForm;
