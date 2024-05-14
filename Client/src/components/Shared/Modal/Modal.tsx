import { IoIosCloseCircle } from "react-icons/io";
import ModalStyled from "./Modal.styled";
import ModalContentBoxStyled from "./ModalContentBox/ModalContentBox.styled";
import CloseButtonStyled from "../CloseButton/CloseButton.styled";
import HeadingStyled from "../Heading/Heading.styled";

interface ModalProps {
  heading: string;
  content: any;
  isVisible: boolean;
  onClose: () => void;
}

const Modal = ({ heading, content, isVisible, onClose }: ModalProps) => {
  return isVisible ? (
    <ModalStyled>
      <ModalContentBoxStyled>
        <CloseButtonStyled onClick={() => onClose()}>
          <IoIosCloseCircle />
        </CloseButtonStyled>
        <HeadingStyled>{heading}</HeadingStyled>
        <section>{content}</section>
      </ModalContentBoxStyled>
    </ModalStyled>
  ) : null;
};

export default Modal;
