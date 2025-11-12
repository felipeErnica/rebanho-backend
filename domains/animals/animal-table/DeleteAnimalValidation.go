package animalTable

import (
	"github.com/felipeErnica/rebanho-backend/apiError"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

const DELETE_OBSERVATION = `
	\nOBS.: A exclusão de animais só é recomendada em caso de erros. Em caso de morte e/ou abate, faça o registro apropriado.
	Lembre-se, apagar um animal do sistema apagará, PERMANENTEMENTE, todas as informções ligadas a ele.
`

func validDelete(db *sqlx.DB, entry DeleteAnimalStruct) *apiError.APIError {

	err := hasChildren(db, entry.Id, entry.UserId)
	if err != nil {
		return err
	}

	err = isCalfInLacation(db, entry.Id, entry.UserId)
	if err != nil {
		return err
	}

	err = isEmbryoDonor(db, entry.Id, entry.UserId)
	if err != nil {
		return err
	}

	err = isReproductionBull(db, entry.Id, entry.UserId)
	if err != nil {
		return err
	}

	err = hasBullPastureEntries(db, entry.Id, entry.UserId)
	if err != nil {
		return err
	}

	if entry.CheckLactation {
		err = hasLactations(db, entry.Id, entry.UserId)
		if err != nil {
			return err
		}
	}

	if entry.CheckInsemination {
		err = hasInseminationEntries(db, entry.Id, entry.UserId)
		if err != nil {
			return err
		}
	}

	if entry.CheckSlaughter {
		err = hasSlaughterEntries(db, entry.Id, entry.UserId)
		if err != nil {
			return err
		}
	}

	if entry.CheckBreeding {
		err = hasBreeding(db, entry.Id, entry.UserId)
		if err != nil {
			return err
		}
	}

	if entry.CheckSlaughter {
		err = isReceiverInTransfer(db, entry.Id, entry.UserId)
		if err != nil {
			return err
		}
	}

	return nil
}

func hasChildren(db *sqlx.DB, id string, userId string) *apiError.APIError {

	query := `
		select exists (
			select 1
			from animals
			where deleted_at is null
				and (mother_id = $1 or father_id = $1)
				and user_id = $2
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.DeleteAPIError("Não é possível deletar animais com filhos registrados no sistema" + DELETE_OBSERVATION)
	}

	return nil
}

func hasLactations(db *sqlx.DB, id string, userId string) *apiError.APIError {

	query := `
		select exists (
			select 1
			from lactations
			where deleted_at is null
				and animal_id = $1
				and user_id = $2
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.DeleteWarningKind(
			"InseminationWarning",
			"Este animal possui lactações registradas. Deseja continuar?"+DELETE_OBSERVATION,
		)
	}

	return nil
}

func hasBreeding(db *sqlx.DB, id string, userId string) *apiError.APIError {

	query := `
		select exists (
			select 1
			from natural_breedings
			where deleted_at is null
				and animal_id = $1
				and user_id = $2
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.DeleteWarningKind(
			"BreedingWarning",
			"Este animal possui coberturas registradas. Deseja continuar?"+DELETE_OBSERVATION,
		)
	}

	return nil
}

func isReceiverInTransfer(db *sqlx.DB, id string, userId string) *apiError.APIError {

	query := `
		select exists (
			select 1
			from embryo_transfer
			where deleted_at is null
				and receiver_id = $1
				and user_id = $2
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.DeleteWarningKind(
			"ReceiverWarning",
			"Esta vaca é uma receptora de embriões. Deseja continuar?"+DELETE_OBSERVATION,
		)
	}

	return nil
}
func hasSlaughterEntries(db *sqlx.DB, id string, userId string) *apiError.APIError {

	query := `
		select exists (
			select 1
			from slaughter_entries
			where deleted_at is null
				and animal_id = $1
				and user_id = $2
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.DeleteWarning("Este animal está presente em um registro de abate. Deseja continuar?" + DELETE_OBSERVATION)
	}

	return nil
}

func hasInseminationEntries(db *sqlx.DB, id string, userId string) *apiError.APIError {

	query := `
		select exists (
			select 1
			from insemination_entries
			where deleted_at is null
				and animal_id = $1
				and user_id = $2
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.DeleteWarning("Este animal já foi inseminado. Deseja continuar?" + DELETE_OBSERVATION)
	}

	return nil
}

func isCalfInLacation(db *sqlx.DB, id string, userId string) *apiError.APIError {

	query := `
		select exists (
			select 1
			from lactations
			where deleted_at is null
				and calf_id = $1
				and user_id = $2
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.DeleteAPIError(`
			Não é possível deletar animais que sejam bezerros em lactações! Altere as lactações da mãe do animal!
		` + DELETE_OBSERVATION)
	}

	return nil
}

func isEmbryoDonor(db *sqlx.DB, id string, userId string) *apiError.APIError {

	query := `
		select exists (
			select 1
			from embryo_transfer
			where donor_id = $1
				and user_id = $2
				and deleted_at is null
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.DeleteAPIError(`
			Esta vaca está listada como doadora de embriões. Altere os registros de transferência antes de concluir!
		` + DELETE_OBSERVATION)
	}

	return nil
}

func isReproductionBull(db *sqlx.DB, id string, userId string) *apiError.APIError {
	query := `
		select exists (
			select 1
			from embryo_transfer
			where bull_id = $1
				and user_id = $2
				and deleted_at is null
		) or exists (
			select 1
			from insemiantion_entries
			where bull_id = $1
				and user_id = $2
				and deleted_at is null
		) or exists (
			select 1
			from natural_mating
			where bull_id = $1
				and user_id = $2
				and deleted_at is null
		) 
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.DeleteAPIError(`
			Não é possível apagar um touro de reprodução. Verifique os registros de inseminação,
			transferência embrionária e cobertura ativa!
		` + DELETE_OBSERVATION)
	}

	return nil
}

func hasBullPastureEntries(db *sqlx.DB, id string, userId string) *apiError.APIError {

	query := `
		select exists (
			select 1
			from pasture_entries pe
				join animals a on a.id = pe.animal_id
			where deleted_at is null
				and animal_type = 'REPRODUCTION_ANIMAL'
				and animal_id = $1
				and user_id = $2
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.DeleteWarning(`
			Não é possível excluir um touro com entrada nos Lotes. Corrija os registros
			de entrada de Lote!
		` + DELETE_OBSERVATION)
	}

	return nil
}
