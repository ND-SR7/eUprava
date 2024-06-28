import { useState } from "react";
import toast from "react-hot-toast";

import Button from "../../Shared/Button/Button";
import Form from "../../Shared/Form/Form";
import Modal from "../../Shared/Modal/Modal";
import HearingCardStyled from "./HearingCard/HearingCard.styled";
import HearingListStyled from "./HearingList.styled";

import CourtHearing from "../../../models/Court/CourtHearing";

import { updateHearing } from "../../../services/CourtService";

interface HearingsProps {
  hearings: CourtHearing[];
  closeParent: () => void;
};

const HearingList = ({hearings, closeParent}: HearingsProps) => {
  const [modalVisible, setModalVisible] = useState(false);
  const [modalContent, setModalContent] = useState<any>(null);
  const [modalHeading, setModalHeading] = useState("");
  const closeModal = () => {
    setModalVisible(false);
  };

  const rescheduleHearing = (formData: any) => {
    const retVal = updateHearing(formData["hearingID"], formData["dateTime"]);
    retVal.then(() => {
      toast.success("Successfully rescheduled court hearing");
      closeModal();
      closeParent();
    }).catch((error) => {
      toast.error("Failed to reschedule court hearing");
      console.error(error);
    });
  };
  
  const newReschedule = (id: string, dateTime: string, reason: string) => {
    dateTime = dateTime.substring(0, dateTime.length - 1);
    const rescheduleFormFields = [
      {attrName: "hearingID", label: "ID", type: "text", value: id, disabled: true},
      {attrName: "dateTime", label: "New date and time", type: "datetime-local", value: dateTime}
    ];

    setModalHeading(`Rescheduling ${reason}`);
    setModalContent(
      <Form 
        heading=""
        formFields={rescheduleFormFields}
        onSubmit={rescheduleHearing} />
    );
    setModalVisible(true);
  };

  const content = hearings.map(hearing =>
    <HearingCardStyled>
      <h1>{hearing.reason}</h1>
      <h6>{hearing.id}</h6>
      <h3>Date and time: {hearing.dateTime.replace("T", " ").replace("Z", "")}</h3>
      <p><b>Court: {hearing.court}</b></p>
      <Button
        buttonType="button"
        label="Reschedule"
        onClick={() => newReschedule(hearing.id!, hearing.dateTime, hearing.reason)} />
    </HearingCardStyled>
  );

  return (
    <>
      <HearingListStyled>{content}</HearingListStyled>
      <Modal
        heading={modalHeading}
        content={modalContent}
        isVisible={modalVisible}
        onClose={closeModal} />
    </>
  );
};

export default HearingList;
