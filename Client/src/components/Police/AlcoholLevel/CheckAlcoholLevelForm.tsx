import React, { useState } from 'react';
import toast from 'react-hot-toast';
import { FormContainer, Input, StyledHeading, SubmitButton } from './CheckAlcoholLevelForm.styled';
import { useNavigate } from 'react-router-dom';
import { checkAlcoholLevel } from '../../../services/PoliceService';

interface CheckAlcoholLevelFormProps {
  closeModal: () => void;
}

const CheckAlcoholLevelForm: React.FC<CheckAlcoholLevelFormProps> = ({ closeModal }) => {
  const [alcoholLevel, setAlcoholLevel] = useState('');
  const [jmbg, setJmbg] = useState('');
  const [location, setLocation] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const data = {
      alcoholLevel: parseFloat(alcoholLevel),
      jmbg,
      location,
    };

    try {
      const response = await checkAlcoholLevel(data);
      toast.success(response.data.message || 'Alcohol level checked successfully');
      navigate('/home/police');
      closeModal();
    } catch (error: any) {
      console.error(error);
      toast.error('Failed to check alcohol level, check driver JMBG');
    }
  };

  return (
    <FormContainer>
      <StyledHeading>Check Alcohol Level</StyledHeading>
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
        <SubmitButton buttonType="submit" label="Check Alcohol Level" />
      </form>
    </FormContainer>
  );
};

export default CheckAlcoholLevelForm;
