import React, { useState } from 'react';
import FormStyled from './Form.styled';
import Button from '../Button/Button';
import Input from '../Input/Input';
import FormBoxStyled from './FormBox/FormBox.styled';

interface FormField {
  type: string;
  label: string;
  attrName: string;
  value?: string | string[];
  options?: string[];
  min?: string;
  max?: string;
  onchange?: (event: React.ChangeEvent<HTMLInputElement>) => void;
}

interface FormProps {
  heading: string;
  formFields: FormField[];
  onSubmit: (formData: { [key: string]: string | string[] }) => void;
}

const Form = ({ heading, formFields, onSubmit }: FormProps) => {
  const initialFormData: { [key: string]: string | string[] } = {};
  formFields.forEach((field) => {
    initialFormData[field.attrName] = field.value ? field.value : '';
    if (field.type === 'checkbox' && !field.value) {
      initialFormData[field.attrName] = [];
    }
  });

  const [formData, setFormData] = useState(initialFormData);

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>, fieldName: string) => {
    let value;
    if (event.target.type === 'checkbox') {
      if (event.target.checked) {
        value = [...formData[fieldName] as string[], event.target.value];
      } else {
        value = (formData[fieldName] as string[]).filter((item: string) => item !== event.target.value);
      }
    } else {
      value = event.target.value;
    }
    setFormData({ ...formData, [fieldName]: value });
  };

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    onSubmit(formData);
    event.preventDefault();
  };

  const showFields = formFields.map((field, index) => {
    if (field.type === 'checkbox') {
      return (
        <FormBoxStyled key={index}>
          <label>{field.label}</label>
          {field.options?.map((option, optionIndex) => (
            <div key={optionIndex}>
              <input
                type="checkbox"
                id={`checkbox_${field.attrName}_${optionIndex}`}
                name={field.attrName}
                value={option}
                checked={formData[field.attrName].includes(option)}
                onChange={(event) => handleChange(event, field.attrName)}
              />
              <label htmlFor={`checkbox_${field.attrName}_${optionIndex}`}>{option}</label>
            </div>
          ))}
        </FormBoxStyled>
      );
    } else if (field.type === 'radio') {
      return (
        <FormBoxStyled key={index}>
          <label>{field.label}</label>
          {field.options?.map((option, optionIndex) => (
            <div key={optionIndex}>
              <input
                type="radio"
                id={`radio_${field.attrName}_${optionIndex}`}
                name={field.attrName}
                value={option}
                checked={formData[field.attrName] === option}
                onChange={(event) => handleChange(event, field.attrName)}
              />
              <label htmlFor={`radio_${field.attrName}_${optionIndex}`}>{option}</label>
            </div>
          ))}
        </FormBoxStyled>
      );
    }  else if (field.type === 'file') {
      return (
        <FormBoxStyled key={index}>
          <Input
            label={field.label}
            id={`input_${index + 1}`}
            attrName={field.attrName}
            type="file"
            accept="image/*"
            handleChange={(event: any):void => {
              if (field.onchange) {
                field.onchange(event);
              }
              handleChange(event, field.attrName);
            }}
            data={""}
          />
        </FormBoxStyled>
      );
    } else {
      return (
        <FormBoxStyled key={index}>
          <Input
            label={field.label}
            id={`input_${index + 1}`}
            attrName={field.attrName}
            type={field.type}
            handleChange={handleChange}
            data={formData[field.attrName]}
            min={field.min}
            max={field.max} />
        </FormBoxStyled>
      );
    }
  });

  return (
    <FormStyled>
      <h1>{heading}</h1>
      <form onSubmit={handleSubmit}>
        {showFields}
        <Button id="submitBtn" buttonType="submit" label="Submit" />
      </form>
    </FormStyled>
  );
};

export default Form;
