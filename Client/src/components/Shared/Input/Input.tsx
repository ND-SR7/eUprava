import LabelStyled from "../Label/Label.styled";
import InputStyled from "./Input.styled";

interface InputProps {
  type: string;
  label: string;
  id: string;
  attrName: string;
  handleChange: ( ...args: any[]) => void;
  data: string | number | boolean | string[];
  accept?: string;
  min?: string;
  max?: string;
  disabled?: boolean;
}

type InputType = string | number;

const Input = ( {type, label, id, attrName, handleChange, data, accept, min, max, disabled} : InputProps ) => {
  function setInputValue(): InputType {
    if (type === "number")
      return data as number;

    return data as string;
  } 

  return (
    <LabelStyled>
      {label}:
        <InputStyled
          id={id}
          name={id}
          type={type}
          value={setInputValue()}
          onChange={(e) => handleChange(e, attrName)}
          accept={accept ? accept : ""}
          min={min ? min : ""}
          max={max ? max : ""}
          disabled={disabled}
        />
    </LabelStyled>
  );
};

export default Input;
